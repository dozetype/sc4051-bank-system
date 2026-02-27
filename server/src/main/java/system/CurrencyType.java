package system;

public enum CurrencyType {
    SGD,
    EUR,
    USD,
    JPY;

    // This helper searches the enum constants for a match
    public static CurrencyType fromString(String text) {
        for (CurrencyType type : CurrencyType.values()) {
            if (type.name().equalsIgnoreCase(text)) {
                return type;
            }
        }
        return null; // Return null if the input doesn't match any currency
    }
}