package system;

import java.net.DatagramPacket;
import java.nio.charset.StandardCharsets;
import java.util.List;

public class AtLeastOnce implements InvocationSemantic{
    private final ParseHandler parse = new ParseHandler();

    @Override
    public String handleRequest(RequestHandler handler, MonitorHandler monitorHandler, DatagramPacket packet){
        String reply;
        String rawRequest = new String(packet.getData(), 0, packet.getLength(), StandardCharsets.UTF_8);
        List<String> parts = parse.parsePacket(rawRequest);

        System.out.println("Received: " + rawRequest + " Request from: "+ packet.getAddress() +" port"+ packet.getPort());
        
        if ("MONITOR".equals(parts.getFirst())){
            reply = monitorHandler.register(parts.get(1), packet.getAddress(), packet.getPort());
        }
        else reply = handler.handleRequest(parts);

        return reply;
    }
}
