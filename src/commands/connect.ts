import shell from "shelljs";

export default async function () {
  try {
    const vectorInstalled = shell.exec("vector --help", {
      silent: true,
    });

    if (vectorInstalled.stderr) {
      throw new Error("vector doesnt exist");
    }

    console.log(vectorInstalled, "vector exists");
  } catch {
    console.log(
      "Vector does not exist, please install at https://vector.dev/docs/setup/installation/"
    );
  }
}
