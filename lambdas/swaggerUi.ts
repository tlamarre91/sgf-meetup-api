import 'dotenv/config';
import type {
	APIGatewayEvent,
	// APIGatewayProxyEventQueryStringParameters,
	Handler,
} from 'aws-lambda';
import { absolutePath } from 'swagger-ui-dist'
// import { readFileSync } from 'fs';

const swaggerJson = {
	testKey: 'testVal'
};

export const handler: Handler = async (event: APIGatewayEvent) => {
	try {
		console.log({ event });

		if (event.path === '/swagger/swagger.json') {
			return {
				statusCode: 200,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(swaggerJson)
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
