package system;

public class Account {
    private final String username; // CASE SENSITIVE
    private final String password;
    private final int accountId;
    private float balance;
    private CurrencyType currency;

    public Account(String username, String password, CurrencyType currency, float balance, int accountId) {
        this.username = username;
        this.password = password;
        this.currency = currency;
        this.balance = balance;
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
    public float getBalance() {
        return balance;
    }
    public CurrencyType getCurrency() {
        return currency;
    }

    public void setCurrency(CurrencyType currency) {
        this.currency = currency;
    }
    public void setBalance(float balance) {
        this.balance += balance;
    }
}
