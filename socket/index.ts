import * as dotenv from "dotenv";
dotenv.config({ path: __dirname + '/.env' });
import "reflect-metadata";
import express from "express";
import { createConnection } from "typeorm";
import http from "http";

const PORT = process.env.PORT || 4000;
const app = express();
const server = http.createServer(app);
const io = require("socket.io")(server, {
    cors: {
        origin: '*',
        methods: ['GET', 'POST'],
        allowEIO3: true
    }
});

const idToSocket: any = {
}
const socketToId:any = {}

// Allow shared fonts.
app.use((_req, res, next) => {
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader("Access-Control-Allow-Headers", "*");
    next();
});

app.set('view engine', 'ejs');
app.use(express.json(), express.urlencoded({ extended: true }));

app.get("/", (_req, res) => {
    res.json({msg: "Hello World"})
});

app.get("/newConnection", async (req, res) => {
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["token"] || !req.query["id"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }

    const { token, id } = req.query;
    idToSocket[token.toString()] = parseInt(id.toString())
    res.status(200).send();
});

app.get("/groupDeleted", async (req, res) => {
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["particapants"] || !req.query["pendingParticapants"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { particapants, pendingParticapants, groupId } = req.query;
    try {
        const particapants2 = JSON.parse(particapants.toString()) as number[];
        for (let i=0;i<particapants2.length;i++){
            const check = socketToId[particapants2[i]];
            if (check) {
                io.to(check).emit("deleted", groupId);
            }
        }
    } catch (err) {
        console.error(err);
    }
    try {
        const pendingParticapants2 = JSON.parse(pendingParticapants.toString()) as number[];
        for (let i=0;i<pendingParticapants2.length;i++){
            const check = socketToId[pendingParticapants2[i]];
            if (check) {
                io.to(check).emit("pendingDeleted", groupId);
            }
        }
    } catch (err) {
        console.error(err);
    }
    console.log(particapants, pendingParticapants, groupId);
});

app.get("/groupAccepted", async (req, res) => {
    console.log(req.headers.authorization);
    console.log(req.query);
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["userId"] || !req.query["owner"] || !req.query["othersCanAdd"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { userId, groupId, owner, othersCanAdd } = req.query;
    try {
        io.to(socketToId[parseInt(userId.toString())]).emit("groupAccepted", {groupId, owner, othersCanAdd:othersCanAdd==="1"});
    } catch (err) {
        console.error(err);
    }
});

app.get("/particapantDeleted", async (req, res) => {
    console.log(req.headers.authorization);
    console.log(req.query);
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["userId"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { userId, groupId } = req.query;
    try {
        io.to(socketToId[parseInt(userId.toString())]).emit("groupRemove", groupId);
    } catch (err) {
        console.error(err);
    }
});

app.get("/particapantDeletedPending", async (req, res) => {
    console.log(req.headers.authorization);
    console.log(req.query);
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["userId"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { userId, groupId } = req.query;
    try {
        io.to(socketToId[parseInt(userId.toString())]).emit("pendingGroupRemove", groupId);
    } catch (err) {
        console.error(err);
    }
});

app.get("/userLeftTransfered", async (req, res) => {
    console.log(req.headers.authorization);
    console.log(req.query);
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["newOwner"] || !req.query["particapants"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { particapants, groupId, newOwner } = req.query;
    try {
        const particapants2 = JSON.parse(particapants.toString()) as number[];
        for (let i=0;i<particapants2.length;i++){
            const check = socketToId[particapants2[i]];
            if (check) {
                io.to(check).emit("newGroupOwner", {groupId, newOwner});
            }
        }
    } catch (err) {
        console.error(err);
    }
});

app.get("/newPendingUser", async (req, res) => {
    console.log(req.headers.authorization);
    console.log(req.query);
    if (!req.headers.authorization) {
        res.status(400).send();
        return
    }

    const header = req.headers.authorization.split(" ")
    if (header.length != 2) {
        res.status(400).send();
        return
    }

    if (!req.query["groupId"] || !req.query["newUser"] || !req.query["ownerId"]) {
        res.status(400).send();
        return
    }

    const secret = header[1];

    if (secret !== process.env.SECRET) {
        res.status(403).send();
        return
    }
    res.status(200).send();
    const { newUser, groupId, ownerId } = req.query;
    try {
        io.to(socketToId[parseInt(ownerId.toString())]).emit("newPendingUser", {groupId, newUser});
    } catch (err) {
        console.error(err);
    }
});

createConnection().then((_data) => {
    console.log("[server] Connection created to database.");
});

io.on('disconnect', (socket: any) => {
    socket.disconnect()
    console.log('🔥: A user disconnected');
});

io.on("connection", (socket: any) => {
    console.log(`Connection from ${socket.id}`);

    socket.on("init", async (data: any) => {
        console.log(`UserId: ${idToSocket[data]}`);
        socketToId[idToSocket[data]] = socket.id;
        io.to(socket.id).emit("data", "Nice bro");
        //io.to(socket.id).emit("blacklist", {success: true, blacklist: data.blackList}); 
    });
});


server.listen(PORT, () => {
    console.log(`[server] Running on port ${PORT}.`);
});