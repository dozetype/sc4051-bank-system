package system;

import java.util.List;
import java.util.Map;

public class RequestHandler {
    private final AccountHandler accHandler = new AccountHandler();
    private final ParseHandler parse = new ParseHandler();

    public String handleRequest(List<String> parts) {
        try {
            if (parts.isEmpty()) return "ERROR: Empty Packet";

            String command = parts.getFirst();

            switch (command) {
                case "CREATEACCOUNT" -> {
                    return createAccount(parts);
                }
                case "CLOSEACCOUNT" -> {
                    return closeAccount(parts);
                }
                case "DEPOSIT" -> {
                    return deposit(parts);
                }
                case "VIEW" -> {
                    return view(parts);
                }
                case "TRANSFER" -> {
                    return transfer(parts);
                }
                default -> {
                }
            }

            String message = "UNKNOWN_COMMAND";
            return message.length() + ": " + message;
        } catch (Exception e) {
            System.err.println("Error processing request: " + e.getMessage());
            return "ERROR: Bad Request";
        }
    }

    private String createAccount(List<String> data) {
        if (data.size() < 5) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        String password = data.get(2);
        CurrencyType currency = CurrencyType.fromString(data.get(3));
        if (currency == null) {
            return "27:FAIL: Invalid currency type";
        }
        float initialBalance = Float.parseFloat(data.get(4));

        Integer accountID = accHandler.getIDCounter();
        Integer newID = accHandler.addAccount(new Account(username, password, currency, initialBalance, accountID));

        return "20:CREATEACCOUNTSUCCESS"+Integer.toString(newID).length()+":"+newID;
    }

    /**
     *
     * @param data
     * (1) Username,
     * (2) ID,
     * (3) Password
     * @return Success or Error Message
     */
    private String closeAccount(List<String> data) {
        if (data.size() < 4) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        Integer accountID = parse.tryParseInt(data.get(2));
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
        Integer accountID = parse.tryParseInt(data.get(2));
        String password = data.get(3);
        CurrencyType currency = CurrencyType.fromString(data.get(4));
        if (currency == null) {
            return "27:FAIL: Invalid currency type";
        }
        float depositAmount = Float.parseFloat(data.get(5));

        Account acc = accHandler.getAccountByID(accountID);
        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        if (acc.getBalance(currency) + depositAmount < 0) return "26:FAIL: Insufficient balance";

        String currBalance = String.valueOf(acc.updateBalance(currency, depositAmount));
        return "14:DEPOSITSUCCESS"+currBalance.length()+":"+currBalance;
    }

    private String view(List<String> data) {
        if (data.size() < 4) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        Integer accountID = parse.tryParseInt(data.get(2));
        String password = data.get(3);

        Account acc = accHandler.getAccountByID(accountID);
        
        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        // Construct account current balance
        StringBuilder reply = new StringBuilder("11:VIEWSUCCESS");
        for (Map.Entry<CurrencyType, Float> entry : acc.getBalances().entrySet()) {
            String curr = entry.getKey().toString();
            String bal = String.valueOf(entry.getValue());
            
            reply.append(curr.length()).append(":").append(curr);
            reply.append(bal.length()).append(":").append(bal);
        }
        
        return reply.toString();
    }

    /**
     * @param data
     * (1) Username,
     * (2) ID,
     * (3) Password,
     * (4) CurrencyType,
     * (5) Amount,
     * (6) Reciever ID
     * @return
     */
    private String transfer(List<String> data) { // MAKE SURE TO NOT TRANSFER NEGATIVE AMOUNT ON CLIENT SIDE
        if (data.size() < 7) return "23:FAIL: Missing arguments";

        // 1. Parse and Basic Validation
        String username = data.get(1);
        Integer accountID = parse.tryParseInt(data.get(2));
        String password = data.get(3);
        CurrencyType currency = CurrencyType.fromString(data.get(4));
        if (currency == null) {
            return "27:FAIL: Invalid currency type";
        }
        float amount = Float.parseFloat(data.get(5));
        Integer recieverID = parse.tryParseInt(data.get(6));

        // 2. Authentication
        Account acc = accHandler.getAccountByID(accountID);
        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        // 3. Receiver Validation
        Account reciever = accHandler.getAccountByID(recieverID);
        if (reciever == null) return "23:FAIL: Account not found";

        // 4. Balance Check
        if (acc.getBalance(currency) - amount < 0) return "26:FAIL: Insufficient balance";

        // 5. Transfering
        acc.updateBalance(currency, -amount);
        reciever.updateBalance(currency, amount);

        // 6. Response Construction
        StringBuilder reply = new StringBuilder("15:TRANSFERSUCCESS");
        for (Map.Entry<CurrencyType, Float> entry : acc.getBalances().entrySet()) {
            String curr = entry.getKey().toString();
            String bal = String.valueOf(entry.getValue());
            
            reply.append(curr.length()).append(":").append(curr);
            reply.append(bal.length()).append(":").append(bal);
        }
        
        return reply.toString();
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
}