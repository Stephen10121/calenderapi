export default function socketConnection(io: any) {
    io.on("connection", (socket: any) => {
        console.log(`Connection from ${socket.id}`);
    
        socket.on("data", async (data: any) => {
            console.log(data);
            //io.to(socket.id).emit("blacklist", {success: true, blacklist: data.blackList}); 
        });
    });
}