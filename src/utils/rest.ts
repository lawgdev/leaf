import axios, { AxiosRequestConfig } from "axios";
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

    const config: AxiosRequestConfig = {
      method,
      headers,
      withCredentials: true,
      data: body ? JSON.stringify(body ?? {}) : undefined,
    };

    const response = await axios(env.API_URL_V1 + endpoint, config);

    if (response.status === 204) {
      return { success: true };
    }

    const responseData: APIResponse<T> = response.data;

    if (response.status >= 300) {
      console.log(responseData);
      return {
        ...responseData,
        error: {
          code: responseData?.error?.code ?? "internal_scoped_error",
          message: responseData?.error?.message ?? "Unknown error",
        },
      };
    }

    return responseData;
  } catch (error) {
    // @ts-ignore todo: type this xD
    if ("response" in error && error.response.data) {
      // Tod: remove this any
      const responseData = (error.response as any).data as APIResponse<T>;

      return {
        success: false,
        error: {
          code: responseData?.error?.code ?? "internal_scoped_error",
          message: responseData?.error?.message ?? "Unknown error",
        },
      };
    }

    return {
      success: false,
      error: {
        code: "internal_scoped_error",
        message: error instanceof Error ? error.message : "Unknown error",
      },
    };
  }
};
