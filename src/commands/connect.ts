import selector from "@inquirer/select";
import { Feed, FeedType, Project } from "@leaf/types";
import { makeRequest } from "@leaf/utils/rest";
import shell from "shelljs";
import {
  amqp,
  awsECS,
  awsS3,
  awsSQS,
  docker,
  file,
  kubernetes,
  redis,
  syslog,
} from "./sources";
import ora from "ora";

const SUPPORTED_SOURCES = {
  docker_logs: { name: "Docker Logs", handler: docker },
  kubernetes_logs: { name: "Kubernetes Logs", handler: kubernetes },
  amqp: { name: "AMQP/RabbitMQ", handler: amqp },
  redis: { name: "Redis", handler: redis },
  aws_ecs_metrics: { name: "AWS ECS Metrics", handler: awsECS },
  aws_s3: { name: "AWS S3", handler: awsS3 },
  aws_sqs: { name: "AWS SQS", handler: awsSQS },
  file: { name: "Log File", handler: file },
  syslog: { name: "Syslog", handler: syslog },
} as const;

export default async function () {
  try {
    const vectorInstalled = shell.exec("vector --help", {
      silent: true,
    });

    if (vectorInstalled.stderr) {
      throw new Error("vector doesn't exist");
    }

    const fetchingProjects = ora("Fetching Projects").start();
    const { success, data, error } = await makeRequest<{ projects: Project[] }>(
      "GET",
      "/users/@me"
    );

    if (!success || !data) {
      fetchingProjects.fail(
        `Failed to fetch projects: ${error?.message ?? "Unknown error"}`
      );

      return;
    }

    fetchingProjects.succeed("Fetched Projects");

    const project = await selector({
      message: "Select Project",
      choices: data.projects.map((p) => ({
        name: p.name,
        value: p.namespace,
        description: `lawg.dev/${p.namespace}`,
      })),
    });

    const selectedProject = data.projects.find(
      (p) => p.namespace === project
    ) as Project;

    const selectableFeeds = selectedProject.feeds
      .filter((f) => f.type === FeedType.APPLICATION)
      .map((f) => ({
        name: f.name,
        value: f.name,
      }));

    if (selectableFeeds.length === 0) {
      console.log(
        "No application feeds found, please create one in our dashboard."
      );
      return;
    }

    const feed = await selector({
      message: "Select An Application Feed",
      choices: selectableFeeds,
    });

    const selectedFeed = selectedProject.feeds.find(
      (f) => f.name === feed
    ) as Feed;

    const source = (await selector({
      message: "Select A Source",
      choices: Object.entries(SUPPORTED_SOURCES).map(([k, v]) => ({
        name: v.name,
        value: k,
      })),
    })) as keyof typeof SUPPORTED_SOURCES;

    SUPPORTED_SOURCES[source].handler(selectedFeed, selectedProject);
  } catch {
    console.log(
      "Vector does not exist, please install at https://vector.dev/docs/setup/installation/"
    );
  }
}
