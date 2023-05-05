import React, { useState, useEffect, useRef } from "react";

const App = () => {
    const [username, setUsername] = useState("");
    const [room, setRoom] = useState("");
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);
    const socketRef = useRef();

    useEffect(() => {
        if (username && room) {
            socketRef.current = new WebSocket(`ws://localhost:8080/ws?username=${username}&room=${room}`);

            socketRef.current.onmessage = (event) => {
                const message = JSON.parse(event.data);
                setMessages((prevMessages) => [...prevMessages, message]);
            };

            socketRef.current.onclose = () => {
                console.log("WebSocket closed");
            };
        }

        return () => {
            if (socketRef.current) {
                socketRef.current.close();
            }
        };
    }, [username, room]);

    const sendMessage = (e) => {
        e.preventDefault();

        if (socketRef.current && message) {
            socketRef.current.send(
                JSON.stringify({
                    username,
                    text: message,
                    room: room,
                })
            );

            setMessage("");
        }
    };

    return (
        <div>
            <h1>WebSocket Chat</h1>
            <form onSubmit={sendMessage}>
                <input
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Username"
                />
                <input
                    type="text"
                    value={room}
                    onChange={(e) => setRoom(e.target.value)}
                    placeholder="Room"
                />
                <input
                    type="text"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    placeholder="Message"
                />
                <button type="submit">Send</button>
            </form>
            <div>
                {messages.map((msg, index) => (
                    <p key={index}>
                        <strong>{msg.username}: </strong>
                        {msg.text}
                    </p>
                ))}
            </div>
        </div>
    );
};

export default App;
