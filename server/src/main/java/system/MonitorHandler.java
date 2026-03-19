package system;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class MonitorHandler {
    private final Map<InetAddress, ClientInfo> clients = new ConcurrentHashMap<>();
    private final ParseHandler parse = new ParseHandler();

    /**
     * @param timeString
     * @param clientIP
     * @param clientPort
     * @return
     */
    public String register(String timeString, InetAddress clientIP, Integer clientPort) {
        Integer duration = parse.tryParseInt(timeString);
        if (duration == null) return "19:FAIL: Invalid Input";

        // Store as one object
        clients.put(clientIP, new ClientInfo(clientPort, System.currentTimeMillis() + (duration * 1000L)));
        return "14:MONITORSUCCESS";
    }

    public void callback(byte[] responseBytes, DatagramSocket socket) {
        long currTime = System.currentTimeMillis();

        // Use an iterator or removeIf to clean up expired clients safely
        clients.entrySet().removeIf(entry -> {
            InetAddress ip = entry.getKey();
            ClientInfo info = entry.getValue();

            if (info.expiryTime < currTime) {
                sendUpdate("14:MONITORTIMESUP".getBytes(), ip, info.port, socket);
                return true; // Removes from map
            } else {
                sendUpdate(responseBytes, ip, info.port, socket);
                return false; // Keeps in map
            }
        });
    }

    private void sendUpdate(byte[] data, InetAddress ip, int port, DatagramSocket socket) {
        try {
            DatagramPacket packet = new DatagramPacket(data, data.length, ip, port);
            socket.send(packet);
        } catch (IOException e) {
            System.err.println("Error sending to " + ip + ":" + port);
        }
    }
}