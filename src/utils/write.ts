import fs from "fs";
import os from "os";

export const writeToPath = async (path: string, text: string) => {
  const homePath = os.homedir();
  const finalPath = `${homePath}/.leaf/${path}`;

  await fs.writeFileSync(finalPath, text);

  return finalPath;
};
