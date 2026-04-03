declare global {
  namespace App {
    interface Locals {
      user?: {
        id: string;
        email: string;
        displayName: string;
      };
      session?: {
        id: string;
        expiresAt: Date;
      };
    }
  }
}

export {};
