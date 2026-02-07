package system;

import java.net.*;
import java.nio.charset.StandardCharsets;

public class Server {
    private static final int PORT = 2222;
    private static final int BUFFER_SIZE = 1000;

    public static void main(String[] args) {
        new Server().start();
    }

    public void start() {
        // We create the "Chef" (Handler) here
        RequestHandler handler = new RequestHandler();

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
                String reply = handler.handleRequest(rawRequest);

                // --- Step 3: Reply (The Waiter brings food back) ---
                byte[] responseBytes = reply.getBytes(StandardCharsets.UTF_8);
                DatagramPacket sendPacket = new DatagramPacket(
                        responseBytes, responseBytes.length,
                        packet.getAddress(), packet.getPort()
                );
                socket.send(sendPacket);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}