package system;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

public class MonitorHandler {
    private final Map<InetAddress, Integer> clientList;
    private final Map<InetAddress, Long> clientExpiry;

    public MonitorHandler() {
        clientList = new HashMap<>();
        clientExpiry = new HashMap<>();
    }

    public String register(String timeString, InetAddress clientIP, Integer clientPort) {
        Integer duration = tryParseInt(timeString);
        if (duration == null) return "19:FAIL: Invalid Input";
        
        clientList.put(clientIP, clientPort);
        clientExpiry.put(clientIP, System.currentTimeMillis() + duration * 1000);
        return "14:MONITORSUCCESS";
    }

    public void callback(byte[] responseBytes, DatagramSocket socket) {
        long currTime = System.currentTimeMillis();
        
        Iterator<Map.Entry<InetAddress, Long>> iterator = clientExpiry.entrySet().iterator();
        while (iterator.hasNext()) {
            Map.Entry<InetAddress, Long> entry = iterator.next();
            if (entry.getValue() < currTime) {
                iterator.remove();
                clientList.remove(entry.getKey());
            }
        }

        for (Map.Entry<InetAddress, Integer> entry : clientList.entrySet()) {
            try {
                DatagramPacket sendPacket = new DatagramPacket(
                    responseBytes, 
                    responseBytes.length,
                    entry.getKey(), 
                    entry.getValue()
                );
                socket.send(sendPacket);
            } catch (IOException e) {
                System.out.println("Failed to send to " + entry.getKey());
            }
        }
    }

    private Integer tryParseInt(String value) {
        try {
            return Integer.valueOf(value);
        } catch (NumberFormatException e) {
            return null;
        }
    }
}