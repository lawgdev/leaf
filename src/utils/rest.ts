import { env } from "@leaf/utils/env";
import { StateManager } from "./stateManager";

export type APIResponse<T = any> = {
  success: boolean;
  data?: T;
  error?: APIError;
};

type CodeErrors =
  | "conflict"
  | "bad_request"
  | "unauthorized"
  | "not_found"
  | "internal_server_error"
  | "forbidden"
  | "internal_scoped_error";

type APIError = {
  code: CodeErrors;
  message: string;
};

/**
 * Make a HTTP request to the Lawg API
 * @param endpoint The API route to call excluding the version discriminator (e.g. auth/login)
 * @param options HTTP Method, request body and other options

 * @returns APIResponse object with the data or error.
 */
export const makeRequest = async <T = any>(
  method: "GET" | "POST" | "PUT" | "PATCH" | "DELETE" | "HEAD",
  endpoint: string,
  body?: any,
  options?: {
    headers?: Record<string, string>;
  }
): Promise<APIResponse<T>> => {
  try {
    const headers: Record<string, string> = {
      Authorization: `Bearer ${(await StateManager.getState()).token}`,
      "Content-Type": "application/json",
      ...options?.headers,
    };

    if (body && (method === "GET" || method === "HEAD")) {
      throw new Error("GET requests cannot have a body");
    }

    const response: APIResponse<T> = await fetch(env.API_URL_V1 + endpoint, {
      method,
      headers,
      credentials: "include",
      body: body ? JSON.stringify(body ?? {}) : null,
    }).then(async (res) =>
      res.status === 204
        ? { success: true }
        : res
            .json()
            .then((json) =>
              res.status >= 300
                ? {
                    ...json,
                    error: {
                      code: json.error.code,
                      message: json.error.message,
                    },
                  }
                : json
            )
            .catch(() =>
              res.status >= 300 ? { success: false } : { success: true }
            )
    );

    return (
      response || {
        success: false,
        error: { code: "internal_scoped_error", message: "Unknown error" },
      }
    );
  } catch (error) {
    console.error(error);

    return {
      success: false,
      error: {
        code: "internal_scoped_error",
        message: error instanceof Error ? error.message : "Unknown error",
      },
    };
  }
};
