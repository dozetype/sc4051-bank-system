package system;

import java.util.List;

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

        Integer accountId = accHandler.getIDCounter();
        Integer newID = accHandler.addAccount(new Account(username, password, currency, initialBalance, accountId));

        return "20:CREATEACCOUNTSUCCESS"+Integer.toString(newID).length()+":"+newID;
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
        Integer accountID = parse.tryParseInt(data.get(2));
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
        Integer accountId = parse.tryParseInt(data.get(2));
        if (accountId == null) return "19:FAIL: Invalid Input";
        String password = data.get(3);
        CurrencyType currency = CurrencyType.fromString(data.get(4));
        if (currency == null) {
            return "27:FAIL: Invalid currency type";
        }
        float depositAmount = Float.parseFloat(data.get(5));

        Account acc = accHandler.getAccountByID(accountId);
        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        String currBalance = String.valueOf(acc.setBalance(depositAmount));
        return "14:DEPOSITSUCCESS"+currBalance.length()+":"+currBalance;
    }

    private String view(List<String> data) {
        if (data.size() < 4) return "23:FAIL: Missing arguments";

        String username = data.get(1);
        Integer accountId = parse.tryParseInt(data.get(2));
        if (accountId == null) return "19:FAIL: Invalid Input";
        String password = data.get(3);

        Account acc = accHandler.getAccountByID(accountId);
        String authError = authenticate(acc, password, username);
        if (authError != null) return authError;

        String currBalance = String.valueOf(acc.getBalance());
        return "11:VIEWSUCCESS"+currBalance.length()+":"+currBalance;
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