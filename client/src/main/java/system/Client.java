package system;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketTimeoutException;

import ui.UI;

public class Client {

    // 1. Settings are moved to the top for easy changing
    private static final String SERVER_IP = "10.91.142.119"; // localhost
    private static final int SERVER_PORT = 2222;
    private static final int TIMEOUT_MS = 3000;
    private static final int BUFFER_SIZE = 512;
    private final UI ui = new UI();
    private String accountName;
    private final char[] accountPassword = new char[12];
    private int accountID;

    public static void main(String[] args) {
        new Client().run();
    }

    private void run() {
        loginOrCreate();
        menu();
        exit();
    }

    private void loginOrCreate() {
        System.out.println("Welcome! (Type 'exit' to quit at any time)");

        boolean isLoggedIn = false;

        try (DatagramSocket socket = new DatagramSocket()) {
            socket.setSoTimeout(TIMEOUT_MS);
            InetAddress serverAddress = InetAddress.getByName(SERVER_IP);

            while (!isLoggedIn) {
                System.out.println("\n--- START MENU ---");
                System.out.println("1. Login");
                System.out.println("2. Create New Account");
                System.out.print("Select: ");

                String choice = ui.inputString();
                if (choice.equalsIgnoreCase("exit")) {
                    exit();
                }

                String requestProtocol = "";

                // 1. Gather Input & Build Protocol String
                if (choice.equals("1")) {
                    System.out.print("Username: ");
                    String user = ui.inputString();
                    System.out.print("Password: ");
                    String pass = ui.inputString();
                    // Format: LOGIN:user:pass
                    requestProtocol = "5:LOGIN" + user.length() + ":" + user + pass.length() + ":" + pass;

                } else if (choice.equals("2")) {
                    System.out.print("New Username: ");
                    String user = ui.inputString();
                    System.out.print("New Password: ");
                    String pass = ui.inputString();
                    System.out.print("Currency: ");
                    String currency = ui.inputString();
                    System.out.print("Initial Deposit: ");
                    String initialDeposit = ui.inputString();
                    // Format: REGISTER:user:pass
                    requestProtocol = "8:REGISTER" + user.length() + ":" + user + pass.length() + ":" + pass + currency.length() + ":" + currency + initialDeposit.length() + ":" + initialDeposit;

                } else {
                    System.out.println("Invalid option.");
                    continue; // Restart loop
                }

                // 2. Send to Server & Check Reply
                String serverReply = sendAndReceiveWithReply(socket, serverAddress, requestProtocol);

                if (serverReply != null) {
                    System.out.println("Reply: " + serverReply);
                }
            }
        } catch (Exception e) {
            System.err.println("Critical Error: " + e.getMessage());
            exit();
        }
    }

    private void menu() {
    }

    private String sendAndReceiveWithReply(DatagramSocket socket, InetAddress address, String message) {
        byte[] requestData = message.getBytes();
        DatagramPacket request = new DatagramPacket(requestData, requestData.length, address, SERVER_PORT);
        byte[] buffer = new byte[BUFFER_SIZE];
        DatagramPacket reply = new DatagramPacket(buffer, buffer.length);

        // 1. Loop for the number of retries
        while (true) {
            try {
                // 2. Send the packet
                socket.send(request);

                // 3. Wait for reply (This blocks until timeout)
                socket.receive(reply);

                // 4. Success! Return the data immediately
                return new String(reply.getData(), 0, reply.getLength());

            } catch (SocketTimeoutException e) {
                // 5. Timeout happened - Log it and let the loop continue
                System.out.println(e.getMessage());
            } catch (IOException e) {
                // Fatal error (e.g., Network Unreachable) - Stop trying
                System.err.println("Network Error: " + e.getMessage());
                return null;
            }
        }
    }

    private void exit() {
        System.out.println("\nThank you for using our application!");
        System.exit(0);
    }
}
