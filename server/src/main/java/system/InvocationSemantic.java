package system;

import java.net.DatagramPacket;

public interface InvocationSemantic {

    String handleRequest(RequestHandler handler, MonitorHandler monitorHandler, DatagramPacket packet);

}
