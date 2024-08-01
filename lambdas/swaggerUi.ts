import 'dotenv/config';
import type {
	APIGatewayEvent,
	// APIGatewayProxyEventQueryStringParameters,
	Handler,
} from 'aws-lambda';
import { absolutePath } from 'swagger-ui-dist'
import { readFileSync } from 'fs';

export const handler: Handler = async (event: APIGatewayEvent) => {
	try {
		if (event.path === '/swagger.json') {
			const swaggerJson = readFileSync('./swagger.json');
			console.log({ swaggerJson });
			return {
				statusCode: 200,
				headers: {
					'Content-Type': 'application/json'
				},
				body: swaggerJson
			};
		}
		console.log({ absolutePath: absolutePath() });

		const body = JSON.stringify({ success: true, event, absolutePath: absolutePath() ?? "something else???" });

		return {
			statusCode: 200,
			body,
			headers: { "Content-Type": "text/html" }
		};
	} catch (err: any) {
		return {
			statusCode: 500,
		};
	}
}
