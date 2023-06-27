import { program } from "commander";
import packageJSON from "../package.json";
import init from "@leaf/commands/init";
import connect from "@leaf/commands/connect";

program
  .name("leaf")
  .description(
    "The easiest way to integrate lawg in your application and have an interface to send logs."
  )
  .version(packageJSON.version);

program
  .command("init")
  .description("Initialize leaf with your lawg token.")
  .action(init);

program
  .command("connect")
  .description("Connect a feed with an application")
  .action(connect);

program.parse();
