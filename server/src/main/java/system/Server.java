package system;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.nio.charset.StandardCharsets;

public class Server {
    private static final int PORT = 2222;
    private static final int BUFFER_SIZE = 1000;
    private InvocationSemantic invocationSemantic;

    public static void main(String[] args) {
        new Server().start();
    }

    public void start() {
        bootUp();
        run();
    }

    private void bootUp(){
        UI ui = new UI();
        System.out.println("\nBOOTING UP...\nSelect Invocation Semantics:\n  1) At-Least-Once\n  2) At-Most-Once");
        while(invocationSemantic == null){
            Integer num = ui.inputInt();
            switch (num) {
                case 1 -> {
                    invocationSemantic = new AtLeastOnce();
                    System.out.println("Running At-Least-Once Invocation Semantics");
                }
                case 2 -> {
                    invocationSemantic = new AtMostOnce();
                    System.out.println("Running At-Most-Once Invocation Semantics");
                }
                default -> System.out.println("Please choose 1 or 2");
            }
        }
    }

    private void run(){
        RequestHandler handler = new RequestHandler();
        MonitorHandler monitorHandler = new MonitorHandler();
        
        try (DatagramSocket socket = new DatagramSocket(PORT)) {
            System.out.println("\n^^^^^^^Server running on port " + PORT + "...^^^^^^^");
            byte[] buffer = new byte[BUFFER_SIZE];

            while (true) {
                // --- Step 1: Receive (The Waiter takes the order) ---
                DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                socket.receive(packet);

                // --- Step 2: Process (The Chef cooks) ---
                String reply;
                reply = invocationSemantic.handleRequest(handler, monitorHandler, packet);

                // --- Step 3: Reply (The Waiter brings food back) ---
                byte[] responseBytes = reply.getBytes(StandardCharsets.UTF_8);
                DatagramPacket sendPacket = new DatagramPacket(
                    reply.getBytes(StandardCharsets.UTF_8),
                    responseBytes.length,
                    packet.getAddress(),
                    packet.getPort()
                );
                socket.send(sendPacket);
                monitorHandler.callback(("8:CALLBACK"+reply).getBytes(StandardCharsets.UTF_8), socket);
            }
        } catch (Exception e) {
            System.err.println(e.getMessage());
        }
    }
}