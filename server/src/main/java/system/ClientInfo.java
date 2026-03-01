package system;

public class ClientInfo {
    int port;
    long expiryTime;

    ClientInfo(int port, long expiryTime) {
        this.port = port;
        this.expiryTime = expiryTime;
    }
}