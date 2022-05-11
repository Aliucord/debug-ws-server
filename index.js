import { spawn } from "child_process";
import readline from "readline";
import { WebSocketServer } from "ws";
import chalk from "chalk";

let connected = false;
const logUtils = {
    incoming: (message) => {
        logUtils.info(chalk.greenBright("<-- ") + message);
    },
    info: (message) => {
        logUtils.clearLine();
        console.info(message);
        if (connected) process.stdout.write(chalk.cyanBright("--> "));
    },
    success: (message) => {
        logUtils.clearLine();
        console.info(chalk.greenBright(message));
    },
    error: (message) => {
        logUtils.clearLine();
        console.error(chalk.redBright(message));
    },
    clearLine: () => {
        process.stdout.clearLine();
        process.stdout.cursorTo(0);
    }
};

const wss = new WebSocketServer({ port: 3000 });

wss.on("connection", async (ws) => {
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    ws.on("message", data => {
        const parsed = JSON.parse(data.toString());
        switch (parsed.level) {
            case 0:
                logUtils.info(`${chalk.bold("T:")} ` + parsed.message);
                break;
            case 1:
                logUtils.info(`${chalk.greenBright("I:")} ` + parsed.message);
                break;
            case 2:
                logUtils.info(`${chalk.yellow("W:")} ` + parsed.message);
                break;
            case 3:
                logUtils.info(`${chalk.redBright("E:")} ` + parsed.message);
                break;
        }
    });
    ws.on("close", () => {
        connected = false;
        rl.close();
        logUtils.error("Websocket connection closed, waiting for reconnection");
    });
    logUtils.incoming("Discord client connected to websocket");

    connected = true;
    while (connected) {
        await new Promise(r => {
            rl.question(chalk.cyanBright("--> "), (cmd) => {
                if (!connected) return;
                else if (["exit", "quit"].includes(cmd)) {
                    ws.close();
                    process.exit();
                } else if (cmd == "clear") {
                    console.clear();
                    process.stdout.write(chalk.cyanBright("--> "));
                } else if (/^\s*$/.test(cmd)) {
                    r();
                } else {
                    ws.send(cmd);
                    r();
                }
            });
        });
    }
});

spawn("adb", ["reverse", "tcp:3000", "tcp:3000"], { stdio: "ignore" }).on("exit", (code) => {
    if (code !== 0) logUtils.error(`Port forwarding port 3000 with adb exited with code ${code}, aliucord may not load`);
    else logUtils.success("Successfully forwarded port 3000 to phone with adb");
});