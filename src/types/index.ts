type PikaIds =
  | "user"
  | "project"
  | "feed"
  | "event"
  | "lawg"
  | "device"
  | "session"
  | "pm";
type Id<T extends PikaIds> = `${T}_${string}`;

type SelfUser = {
  id: Id<"user">;
  username: string;
  email: string;
  github: number;
  api_token: Id<"lawg">;
  icon: string | null;
};

type User = {
  id: Id<"user">;
  username: string;
  icon: string | null;
};

type Feed = {
  id: Id<"feed">;
  project_id: Id<"project">;
  name: string;
  description: string | null;
  emoji: string | null;
  type: FeedType;
};

type Project = {
  id: Id<"project">;
  namespace: string;
  name: string;
  tier: ProjectTier;
  icon: string | null;
  feeds: Feed[];
  members: User[];
};

type Event = {
  id: Id<"event">;
  project_id: Id<"project">;
  feed_id: Id<"feed">;
  title: string;
  description: string | null;
  emoji: string | null;
};

export enum FeedType {
  EVENT,
  APPLICATION,
}

export enum ProjectTier {
  FREE,
  PRO,
  STARTUP,
  ENTERPRISE,
}

export type { SelfUser, User, Feed, Project, Event };
