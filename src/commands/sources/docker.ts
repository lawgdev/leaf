import Docker from "dockerode";
import select from "@inquirer/select";
import shell from "shelljs";
import fs from "fs";

import { Source } from ".";
import { StateManager } from "../../utils/stateManager";
import { writeToPath } from "../../utils/write";
import ora from "ora";

const generateConfig = (
  source: string,
  includeContainer: string,
  project: string,
  feed: string,
  token: string
) => `
${source.trim()}
include_containers = [ "${includeContainer}" ]

[sinks.lawg_sink]
type = "http"
encoding.codec = "json"
inputs = ["source0"]
uri = "http://100.127.114.55:8080/v1/projects/${project}/feeds/${feed}/logs"
auth.strategy = "bearer"
auth.token = "${token}"
`;

export const docker: Source = async (feed, project) => {
  const docker = new Docker();
  const fetchContainers = ora("Fetching containers").start();
  const containers = await docker.listContainers();
  fetchContainers.succeed("Fetched containers");

  const selectedContainer = await select({
    message: "Select a container",
    choices: containers.map((container) => ({
      name: `${container.Names[0].substring(1)} (${container.Id.substring(
        0,
        12
      )})`,
      value: container.Names[0].substring(1),
    })),
  });

  let vectorGeneratedConfig = shell.exec("vector generate docker_logs", {
    silent: true,
  }).stdout;

  // remove first 2 lines of generatedConfig (irrelevant to config) (im waiting for the day this is added https://tc39.es/proposal-pipeline-operator/)
  const lines = vectorGeneratedConfig.split("\n");
  lines.splice(0, 2)[0];

  vectorGeneratedConfig = lines.join("\n");

  const config = generateConfig(
    vectorGeneratedConfig,
    selectedContainer,
    project.namespace,
    feed.name,
    (await StateManager.getState()).token
  );

  const configPath = await writeToPath(
    `configs/${feed.name}-${selectedContainer}.toml`,
    config
  );
  console.log(`Config generated and saved to ${configPath}`);
};
