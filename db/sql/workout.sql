CREATE TABLE "exercises" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" text NOT NULL,
  "description" text,
  "sets" int NOT NULL,
  "href" text,
  "equipment" varchar,
  "target_muscle" varchar NOT NULL,
  "replacement" BIGSERIAL
);

CREATE TABLE "sets" (
  "id" BIGSERIAL PRIMARY KEY,
  "exercise_id" BIGSERIAL NOT NULL,
  "reps" int NOT NULL
);

CREATE INDEX ON "exercises" ("name");

CREATE INDEX ON "exercises" ("equipment");

CREATE INDEX ON "exercises" ("target_muscle");

ALTER TABLE "sets" ADD CONSTRAINT "sets" FOREIGN KEY ("exercise_id") REFERENCES "exercises" ("id");

ALTER TABLE "exercises" ADD FOREIGN KEY ("replacement") REFERENCES "exercises" ("id");
