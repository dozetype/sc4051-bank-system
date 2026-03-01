package system;

import java.net.DatagramPacket;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class AtMostOnce implements InvocationSemantic{
    private final Map<String, String> processedRequest = new ConcurrentHashMap<>();
    private final ParseHandler parse = new ParseHandler();

    /**
     * len(UUID):<UUID>len(COMMAND):<COMMAND>.....
     */
    @Override
    public String handleRequest(RequestHandler handler, MonitorHandler monitorHandler, DatagramPacket packet){
        String reply;
        String rawRequest = new String(packet.getData(), 0, packet.getLength(), StandardCharsets.UTF_8);
        List<String> parts = parse.parsePacket(rawRequest);

        System.out.println("Received: " + rawRequest + " from: "+ packet.getAddress() +" port"+ packet.getPort());

        String uuid = parts.get(0);
        if (processedRequest.containsKey(uuid)) {
            System.out.println("Duplicate detected for UUID: " + uuid);
            return processedRequest.get(uuid);
        }
        else if ("MONITOR".equals(parts.get(1))){
            reply = monitorHandler.register(parts.get(2), packet.getAddress(), packet.getPort());
        }
        else reply = handler.handleRequest(parts);

        processedRequest.put(parts.getFirst(), reply); // Map UUID to Reply
        return reply;
    }
}
