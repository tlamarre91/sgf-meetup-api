import 'dotenv/config';

import { DynamoDBClient } from '@aws-sdk/client-dynamodb';

const isDev = process.env.NODE_ENV === 'development';

function makeDynamoDBClient() {
	const accessKeyId = process.env.LAMBDA_AWS_ACCESS_KEY_ID;
	const secretAccessKey = process.env.LAMBDA_AWS_SECRET_ACCESS_KEY;

	let client: DynamoDBClient;

	if (isDev) {
		// App is running in development mode; return a client that expects DynamoDB to be running locally
		const endpoint = 'http://dynamodb-local:8000'; // TODO: make local endpoint configurable
		const credentials = {
			accessKeyId: accessKeyId ?? "dev-key",
			secretAccessKey: secretAccessKey ?? "dev-key"
		};

		console.log({ credentials });

		client = new DynamoDBClient({ endpoint, credentials });
	} else {
		client = new DynamoDBClient();
	}

	return client;
}

export const dynamoDbClient = makeDynamoDBClient();
