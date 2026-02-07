package system;

import java.util.ArrayList;
import java.util.List;

class RequestHandler {
    private final AccountHandler accHandler = new AccountHandler();

    public String handleRequest(String rawData) {
        try {
            List<String> parts = parsePacket(rawData);

            if (parts.isEmpty()) return "ERROR: Empty Packet";

            String command = parts.getFirst();

            // Simple switch to route commands
            if (command.equals("LOGIN")) {
                return login(parts);
            }
            else if (command.equals("REGISTER")) {
                return register(parts);
            }
            String message = "UNKNOWN_COMMAND";
            return message.length() + ": " + message;
        } catch (Exception e) {
            System.err.println("Error processing request: " + e.getMessage());
            return "ERROR: Bad Request";
        }
    }

    private String login(List<String> data) {
        if (data.size() < 3) return "FAIL: Missing arguments";

        String username = data.get(1);
        String password = data.get(2);

        if (accHandler.getAccountByUsername(username) == null) {
            String message = "FAIL: User not found";
            return message.length() + ":" + message;
        }
        Account acc = accHandler.getAccountByUsername(username);
        if (acc.getPassword().equals(password)) {
            int accountId = acc.getAccountId();
            return "12:LOGINSUCCESS"+Integer.toString(accountId).length()+":"+accountId;
        }
        return "20:FAIL: Wrong password";
    }

    private String register(List<String> data) {
        if (data.size() < 4) return "23:FAIL: Missing arguments";
        String username = data.get(1);
        if (accHandler.getAccountByUsername(username) != null) {
            String message = "FAIL: Username already exists";
            return message.length() + ":" + message;
        }
        String password = data.get(2);
        CurrencyType currency = CurrencyType.valueOf(data.get(3));
        float initialBalance = Float.parseFloat(data.get(4));
        Integer accountId = accHandler.getIDCounter();
        Integer newID = accHandler.addAccount(new Account(username, password, currency, initialBalance, accountId));
        return "15:REGISTERSUCCESS"+Integer.toString(newID).length()+":"+newID;
    }

    // Helper to decode the "Length:Value" format
    private List<String> parsePacket(String packet) {
        List<String> parts = new ArrayList<>();
        int cursor = 0;

        while (cursor < packet.length()) {
            int delimiterPos = packet.indexOf(':', cursor);
            if (delimiterPos == -1) break;

            String lengthStr = packet.substring(cursor, delimiterPos);
            int length = Integer.parseInt(lengthStr);

            int dataStart = delimiterPos + 1;
            String data = packet.substring(dataStart, dataStart + length);

            parts.add(data);
            cursor = dataStart + length;
        }
        return parts;
    }
}