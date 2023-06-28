import { Feed, Project } from "@leaf/types";

export * from "./amqp";
export * from "./aws_ecs";
export * from "./aws_s3";
export * from "./aws_sqs";
export * from "./docker";
export * from "./file";
export * from "./kubernetes";
export * from "./redis";
export * from "./syslog";

export type Source = (feed: Feed, project: Project) => Promise<void>;
