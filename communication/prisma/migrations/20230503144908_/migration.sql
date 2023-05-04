-- CreateTable
CREATE TABLE "players" (
    "id" VARCHAR(28) NOT NULL,
    "username" VARCHAR(16) NOT NULL,
    "status_id" SMALLINT NOT NULL,

    CONSTRAINT "players_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "friend_relationships" (
    "inviter_id" TEXT NOT NULL,
    "acceptor_id" TEXT NOT NULL
);

-- CreateTable
CREATE TABLE "player_statuses" (
    "id" SMALLSERIAL NOT NULL,
    "name" VARCHAR(16) NOT NULL,

    CONSTRAINT "player_statuses_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "players_username_key" ON "players"("username");

-- CreateIndex
CREATE UNIQUE INDEX "friend_relationships_inviter_id_acceptor_id_key" ON "friend_relationships"("inviter_id", "acceptor_id");

-- CreateIndex
CREATE UNIQUE INDEX "player_statuses_name_key" ON "player_statuses"("name");

-- AddForeignKey
ALTER TABLE "players" ADD CONSTRAINT "players_status_id_fkey" FOREIGN KEY ("status_id") REFERENCES "player_statuses"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "friend_relationships" ADD CONSTRAINT "friend_relationships_inviter_id_fkey" FOREIGN KEY ("inviter_id") REFERENCES "players"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "friend_relationships" ADD CONSTRAINT "friend_relationships_acceptor_id_fkey" FOREIGN KEY ("acceptor_id") REFERENCES "players"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
