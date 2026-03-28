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

    /**
     * Registers a new account in the system with an initial balance.
     * @param data A list containing: ["CREATEACCOUNT", Username, Password, CurrencyType, InitialBalance]
     * @return A formatted string "CREATEACCOUNTSUCCESS" followed by the newly generated Account ID.
     */
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
     * @param data A list containing: ["CLOSEACCOUNT", Username, AccountID, Password]
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

    /**
     * Deposits a specified amount of a specific currency into an account.
     * @param data A list containing: ["DEPOSIT", Username, AccountID, Password, CurrencyType, Amount]
     * @return A formatted string "DEPOSITSUCCESS" followed by the new balance of that currency, 
     * or an error message if authentication or validation fails.
     */
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

    /**
     * Retrieves and formats all currency balances for a specific account.
     * @param data A list containing: ["VIEW", Username, AccountID, Password]
     * @return A string starting with "VIEWSUCCESS" followed by a list of 
     * formatted currency names and their respective balances.
     */
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
     * @param data A list containing: ["TRANSFER", Username, AccountID, Password, 
     *                                  CurrencyType, Amount, Reciever ID]
     * @return A string indicating failure (e.g., "FAIL: reasons") or success 
     *         followed by the sender's updated balance list.
     */
    private String transfer(List<String> data) {
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
        if (amount < 0) {
            return "36:FAIL: No transfering negative amount";
        }
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
     * Verifies account credentials before allowing an operation.
     * @param acc The account object retrieved from the database.
     * @param password The password string provided by the user.
     * @param username The username provided (optional, checked only if not null).
     * @return null if authentication succeeds; an error string if it fails.
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