import * as fs from "fs";
import * as os from "os";
import * as path from "path";
import { Project } from "../types";

interface State {
  token: string;
  applications: {
    name: string;
    connected: boolean;
  }[];
}

export class StateManager {
  public static getStatePath(): string {
    const homePath = os.homedir();
    const statePath = path.join(homePath, ".leaf", "state.txt");

    // Create the directory if it doesn't exist
    if (!fs.existsSync(path.dirname(statePath))) {
      fs.mkdirSync(path.dirname(statePath), { recursive: true });
    }

    // Create the file if it doesn't exist
    if (!fs.existsSync(statePath)) {
      fs.writeFileSync(statePath, JSON.stringify({}));
    }

    return statePath;
  }

  public static async getState(): Promise<State> {
    const tokenFilePath = StateManager.getStatePath();
    const content = await fs.readFileSync(tokenFilePath, "utf8");

    try {
      const parsedContent = JSON.parse(content);

      return parsedContent as State;
    } catch {
      return {
        token: "",
        applications: [],
      };
    }
  }

  public static async setState(state: Partial<State>): Promise<void> {
    const tokenFilePath = StateManager.getStatePath();
    const currentState = await StateManager.getState();

    await fs.writeFileSync(
      tokenFilePath,
      JSON.stringify({ ...currentState, ...state })
    );
  }

  public static async addApplication(name: string): Promise<void> {
    const currentState = await StateManager.getState();

    await StateManager.setState({
      applications: [
        ...currentState.applications,
        {
          name,
          connected: false,
        },
      ],
    });
  }
}
