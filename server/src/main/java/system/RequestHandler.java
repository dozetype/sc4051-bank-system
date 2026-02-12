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
        if (data.size() < 3) return "23:FAIL: Missing arguments";

        Integer accountId = tryParseInt(data.get(1));
        if (accountId == null) return "19:FAIL: Invalid Input";
        String password = data.get(2);
        Account acc = accHandler.getAccountByID(accountId);

        String authError = authenticate(acc, password, null);
        if (authError != null) return authError;

        return "12:LOGINSUCCESS"+Integer.toString(accountId).length()+":"+accountId;
    }

    private String register(List<String> data) {
        if (data.size() < 5) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        String password = data.get(2);
        CurrencyType currency = CurrencyType.valueOf(data.get(3)); // TODO: Handle incorrect enum
        float initialBalance = Float.parseFloat(data.get(4));

        Integer accountId = accHandler.getIDCounter();
        Integer newID = accHandler.addAccount(new Account(username, password, currency, initialBalance, accountId));

        return "15:REGISTERSUCCESS"+Integer.toString(newID).length()+":"+newID;
    }

    /**
     *
     * @param data
     * (1)Username
     * (2)ID
     * (3)Password
     * @return Success or Error Message
     */
    private String closeAccount(List<String> data) {
        if (data.size() < 4) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        Integer accountID = tryParseInt(data.get(2));
        if (accountID == null) return "19:FAIL: Invalid Input";
        String password = data.get(3);
        Account acc = accHandler.getAccountByID(accountID);

        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        accHandler.closeAccount(acc);
        return "12:CLOSESUCCESS";
    }

    private String deposit(List<String> data) {
        if (data.size() < 6) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        Integer accountId = tryParseInt(data.get(2));
        if (accountId == null) return "19:FAIL: Invalid Input";
        String password = data.get(3);
        CurrencyType currency = CurrencyType.valueOf(data.get(4)); // TODO: handle Currency
        float depositAmount = Float.parseFloat(data.get(5));
        Account acc = accHandler.getAccountByID(accountId);

        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        String currBalance = String.valueOf(acc.setBalance(depositAmount));
        return "14:DEPOSITSUCCESS"+currBalance.length()+":"+currBalance.length();
    }

    /**
     *
     * @param acc
     * @param password
     * @param username Add Username only when it needs to check against it with acc
     * @return
     */
    private String authenticate(Account acc, String password, String username) {
        if (acc == null) return "23:FAIL: Account not found";
        if (!acc.getPassword().equals(password)) return "20:FAIL: Wrong password";
        if (username != null && !acc.getUsername().equals(username)) {
            return "23:FAIL: Username mismatch";
        }
        return null;
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

    private Integer tryParseInt(String value) {
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            return null; // Instead of crashing, we return null to handle it gracefully
        }
    }
}