import express, { Request, Response, NextFunction } from 'express';
import { handle } from 'gohandler';

const app = express();
const port = process.env.PORT || 3000;

app.get('/nodejs', async (req: Request, res: Response) => {
    res.contentType('application/json');
    await new Promise((resolve) => setTimeout(resolve, 1));
    res.status(200).send(JSON.stringify({ message: "Hello from NodeJS" }));
});

app.use(async (req: Request, res: Response) => {
    const method = req.method;
    const url = req.url;
    const headersMap: Record<string, string[]> = {};
    for (const [key, value] of Object.entries(req.headers)) {
        if (typeof (value) === 'string') {
            headersMap[key] = [value];
        } else if (Array.isArray(value)) {
            headersMap[key] = value;
        }
    }
    const headers = JSON.stringify(headersMap);
    const parts: Buffer[] = [];
    let body: Buffer;

    req.on('data', (chunk: Buffer) => {
        parts.push(chunk);
    });
    await new Promise((resolve) => req.on('end', resolve));

    body = Buffer.concat(parts);

    const { statusCode, headersJson, responseBody } = await handle(method, url, headers, body);

    const responseHeaders = JSON.parse(headersJson) as Record<string, string[]>;
    const responseHeadersMap = new Map<string, string[]>(Object.entries(responseHeaders));

    res.setHeaders(responseHeadersMap).status(statusCode).end(responseBody);
});

const server = app.listen(port, () => {
    console.log(`Server running on port ${port}`);
});

process.on('SIGINT', function () {
    console.log('Do something useful here.');
    server.close();
});
