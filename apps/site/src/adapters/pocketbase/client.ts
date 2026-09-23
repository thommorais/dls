import PocketBase from 'pocketbase'
import { ENVS } from '_/envs'

let client: PocketBase | undefined

export const getPocketBaseClient = (): PocketBase => {
	if (!client) {
		client = new PocketBase(ENVS.API_URL)

		// Vite's dev refresh re-runs effects, and autocancellation turns the
		// second identical request into a rejection that reads as a real
		// failure.
		if (ENVS.IS_DEV) client.autoCancellation(false)
	}

	return client
}
