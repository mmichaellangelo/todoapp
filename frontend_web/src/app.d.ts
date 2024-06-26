// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		interface Error {
			message: string;
			errorId: string;
		}
		interface Locals {
			username?: string;
			userid?: number;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
