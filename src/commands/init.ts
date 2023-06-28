import { input } from "@inquirer/prompts";
import { Project, SelfUser } from "@leaf/types";
import { makeRequest } from "@leaf/utils/rest";
import { StateManager } from "@leaf/utils/stateManager";
import ora from "ora";

export default async function () {
  const token = await input({
    message: "Enter your lawg.dev token (https://app.lawg.dev/user/settings):",
  });

  const fetchingUserOra = ora("Authorizing with lawg.dev").start();
  const { data, success, error } = await makeRequest<{
    user: SelfUser;
    projects: Project[];
  }>("GET", "/users/@me", undefined, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!success || !data) {
    fetchingUserOra.fail(
      `Failed to authorize with lawg.dev: ${
        error?.code === "unauthorized"
          ? "Invalid token"
          : error?.message ?? "An unknown error occurred"
      }`
    );

    return;
  }

  fetchingUserOra.succeed("Authorized with lawg.dev");
  StateManager.setState({ token });

  console.log(`Logged in as ${data.user.username} (${data.user.email})`);
}
