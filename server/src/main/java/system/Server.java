package system;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

public class Server {
    private static final int PORT = 2222;
    private static final int BUFFER_SIZE = 1000;

    public static void main(String[] args) {
        new Server().start();
    }

    public void start() {
        // We create the "Chef" (Handler) here
        RequestHandler handler = new RequestHandler();
        MonitorHandler monitorHandler = new MonitorHandler();

        try (DatagramSocket socket = new DatagramSocket(PORT)) {
            System.out.println("\n^^^^^^^Server running on port " + PORT + "...^^^^^^^");
            byte[] buffer = new byte[BUFFER_SIZE];

            while (true) {
                // --- Step 1: Receive (The Waiter takes the order) ---
                DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                socket.receive(packet);

                String rawRequest = new String(
                        packet.getData(), 0, packet.getLength(), StandardCharsets.UTF_8
                );
                System.out.println("Received: " + rawRequest);

                // --- Step 2: Process (The Chef cooks) ---
                List<String> parts = parsePacket(rawRequest);
                String reply;
                if ("MONITOR".equals(parts.getFirst())){
                    reply = monitorHandler.register(parts.get(1), packet.getAddress(), packet.getPort());
                }
                else reply = handler.handleRequest(parts);

                // --- Step 3: Reply (The Waiter brings food back) ---
                byte[] responseBytes = reply.getBytes(StandardCharsets.UTF_8);
                DatagramPacket sendPacket = new DatagramPacket(
                    responseBytes,
                    responseBytes.length,
                    packet.getAddress(),
                    packet.getPort()
                );
                socket.send(sendPacket);
                monitorHandler.callback(responseBytes, socket);
            }
        } catch (Exception e) {
            System.err.println(e.getMessage());
        }
    }

    // Helper to decode the "Length:Value" format
    private List<String> parsePacket(String packet) {
        List<String> parts = new ArrayList<>();
        int cursor = 0;

        while (cursor < packet.length()) {
            int delimiterPos = packet.indexOf(':', cursor);
            if (delimiterPos == -1) break;

            String lengthStr = packet.substring(cursor, delimiterPos);
            int length = Integer.parseInt(lengthStr);

            int dataStart = delimiterPos + 1;
            String data = packet.substring(dataStart, dataStart + length);

            parts.add(data);
            cursor = dataStart + length;
        }
        return parts;
    }
}