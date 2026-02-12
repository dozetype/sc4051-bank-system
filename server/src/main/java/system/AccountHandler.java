package system;

import java.util.HashMap;
import java.util.Map;

public class AccountHandler {
    private final Map<Integer, Account> accounts = new HashMap<>();
    private Integer IDCounter = 1000;

    public AccountHandler() {
        accounts.put(IDCounter, new Account("tom", "123", CurrencyType.SGD, 50, IDCounter++)); // ID: 1000
        accounts.put(IDCounter, new Account("dick", "123", CurrencyType.USD, 100, IDCounter++)); //ID: 1001
        accounts.put(IDCounter, new Account("harry", "123", CurrencyType.EUR, 290, IDCounter++)); //ID: 1002
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
    public Integer addAccount(Account account) {
        accounts.put(getIDCounter(), account);
        return incrementID();
    }
    public void closeAccount(Account account) {
        accounts.remove(account.getAccountId());
    }
}
