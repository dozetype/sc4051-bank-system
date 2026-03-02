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
    private final double packetLossRate;

    public MonitorHandler(double packetLossRate){
        this.packetLossRate = packetLossRate;
    }

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

        clients.forEach((ip, info) -> {
            if (info.expiryTime < currTime) {
                clients.remove(ip);
            } else {
                try {
                    DatagramPacket sendPacket = new DatagramPacket(responseBytes, responseBytes.length, ip, info.port);
                    if (Math.random() > packetLossRate) socket.send(sendPacket);
                } catch (IOException e) {
                    System.out.println("Failed to send to " + ip);
                }
            }
        });
    }
}