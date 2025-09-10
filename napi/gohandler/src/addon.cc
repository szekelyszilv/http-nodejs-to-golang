#include "napi.h"
#include "golib/gohandler.h"

namespace
{
    class GoServeHTTPWorker : public Napi::AsyncWorker
    {
    private:
        uintptr_t handle;
        Napi::Promise::Deferred deferred;

    public:
        GoServeHTTPWorker(Napi::Env env, uintptr_t handle, Napi::Promise::Deferred const &d)
            : AsyncWorker(env), handle(handle), deferred(d)
        {
        }

        virtual ~GoServeHTTPWorker() override
        {
            GoHandlerFree(handle);
            handle = (uintptr_t)nullptr;
        }

        void Execute() override
        {
            GoHandlerRun(handle);
        }

        void OnOK() override
        {
            Napi::Env env = Env();

            int statusCode = GoHandlerGetResponseStatusCode(handle);

            std::unique_ptr<char, decltype(&free)> responseHeadersPtr(GoHandlerGetResponseHeaders(handle), free);
            Napi::String responseHeaders = Napi::String::New(env, responseHeadersPtr.get());

            size_t responseBodySize = GoHandlerGetResponseBodySize(handle);
            Napi::Buffer<uint8_t> responseBody = Napi::Buffer<uint8_t>::New(env, responseBodySize);
            GoHandlerGetResponseBodyBytes(handle, responseBody.Data(), responseBodySize);

            Napi::Object result = Napi::Object::New(env);
            result.Set(Napi::String::New(env, "statusCode"), statusCode);
            result.Set(Napi::String::New(env, "headersJson"), responseHeaders);
            result.Set(Napi::String::New(env, "responseBody"), responseBody);

            deferred.Resolve(result);
        }
    };

    Napi::Object Handle(const Napi::CallbackInfo &info)
    {
        Napi::Env env = info.Env();

        if (info.Length() < 4)
        {
            Napi::TypeError::New(env, "Expected arguments [method: string, url: string, headers: string, body: buffer]").ThrowAsJavaScriptException();
            return {};
        }

        if (!info[0].IsString())
        {
            Napi::TypeError::New(env, "Expected argument 1 to be a string").ThrowAsJavaScriptException();
            return {};
        }

        if (!info[1].IsString())
        {
            Napi::TypeError::New(env, "Expected argument 2 to be a string").ThrowAsJavaScriptException();
            return {};
        }

        if (!info[2].IsString())
        {
            Napi::TypeError::New(env, "Expected argument 3 to be a string").ThrowAsJavaScriptException();
            return {};
        }

        if (!info[3].IsBuffer())
        {
            Napi::TypeError::New(env, "Expected argument 4 to be a buffer").ThrowAsJavaScriptException();
            return {};
        }

        std::string method = info[0].As<Napi::String>().Utf8Value();
        std::string url = info[1].As<Napi::String>().Utf8Value();
        std::string headers = info[2].As<Napi::String>().Utf8Value();
        Napi::Buffer<uint8_t> body = info[3].As<Napi::Buffer<uint8_t>>();

        char *errStringPtr = nullptr;

        uintptr_t handle = GoHandlerNew(
            const_cast<char *>(method.c_str()),
            const_cast<char *>(url.c_str()),
            const_cast<char *>(headers.c_str()),
            body.Length(),
            body.Data(),
            &errStringPtr);
        if (!handle)
        {
            Napi::String errString = Napi::String::New(env, errStringPtr);
            free(errStringPtr);

            Napi::TypeError::New(env, errString)
                .ThrowAsJavaScriptException();
            return Napi::Object::New(env);
        }

        Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
        (new GoServeHTTPWorker(env, handle, deferred))->Queue();

        return deferred.Promise();
    }

    Napi::Object Init(Napi::Env env, Napi::Object exports)
    {
        exports.Set(Napi::String::New(env, "handle"), Napi::Function::New(env, Handle));
        return exports;
    }
}

NODE_API_MODULE(interop, Init)
