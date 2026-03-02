package system;

import java.util.HashMap;
import java.util.Map;

public class Account {
    private final String username; // CASE SENSITIVE
    private final String password;
    private final int accountId;
    private final Map<CurrencyType, Float> balances;

    public Account(String username, String password, CurrencyType currency, float balance, int accountId) {
        this.username = username;
        this.password = password;
        this.balances = new HashMap<>();
        this.balances.put(currency, balance);
        this.accountId = accountId;
    }

    public String getUsername() {
        return username;
    }
    public String getPassword() {
        return password;
    }
    public int getAccountId() {
        return accountId;
    }
    public Map<CurrencyType, Float> getBalances() {
        return balances;
    }
    public float getBalance(CurrencyType currency) {
        return balances.getOrDefault(currency, 0.0f);
    }

    /**
     * Updates the balance for a specific currency.
     * @param currency The type of currency to update
     * @param amount The amount to add (positive) or subtract (negative)
     * @return The new total balance for that currency
     */
    public float updateBalance(CurrencyType currency, float amount) {
        float currentBalance = balances.getOrDefault(currency, 0.0f);
        float newBalance = currentBalance + amount;
        
        balances.put(currency, newBalance);
        return newBalance;
    }
}
