declare module 'gohandler' {
    export function handle(method: string, url: string, headers: string, body: Buffer): Promise<{ statusCode: number, headersJson: string, responseBody: Buffer }>;
}
