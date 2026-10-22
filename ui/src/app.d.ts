// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    // interface Locals {}
    // interface PageData {}
    interface PageState {
      clientWorkspaceUrl?: string;
      flash?: import("$lib/flash").Flash;
    }
    // interface Platform {}
  }
}

export {};
