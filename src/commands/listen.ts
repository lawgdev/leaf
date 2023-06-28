import shell from "shelljs";
import os from "os";
import ora from "ora";

export default function () {
  const homeDirectory = os.homedir();

  const listeningToLogs = ora("Listening to logs").start();
  shell.exec(`vector -c ${homeDirectory}/.leaf/configs/*.toml`);

  listeningToLogs.succeed("Listening to logs");
}
