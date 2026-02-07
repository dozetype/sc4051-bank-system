package system;

import java.util.HashMap;
import java.util.Map;

public class AccountHandler {
    private final Map<Integer, Account> accounts = new HashMap<>();
    private final Map<String, Integer> usernames = new HashMap<>();
    private Integer IDCounter = 1000;

    public AccountHandler() {
        accounts.put(IDCounter, new Account("tom", "123", CurrencyType.SGD, 50, IDCounter));
        usernames.put("tom", IDCounter++);
        accounts.put(IDCounter, new Account("dick", "123", CurrencyType.USD, 100, IDCounter));
        usernames.put("dick", IDCounter++);
        accounts.put(IDCounter, new Account("harry", "123", CurrencyType.EUR, 290, IDCounter));
        usernames.put("harry", IDCounter++);
        accounts.get(1001).setBalance(50);
    }

    public Integer getIDCounter() {
        return IDCounter;
    }
    public Integer incrementID() {
        return this.IDCounter++;
    }
    public Account getAccountByID(Integer ID) {
        return accounts.get(ID);
    }
    public Account getAccountByUsername(String username) {
        return getAccountByID(usernames.get(username));
    }
    public Integer addAccount(Account account) {
        accounts.put(getIDCounter(), account);
        usernames.put(account.getUsername(), getIDCounter());
        return incrementID();
    }
    public void closeAccount(Account account) {
        usernames.remove(account.getUsername());
        accounts.remove(account.getAccountId());
    }
}
