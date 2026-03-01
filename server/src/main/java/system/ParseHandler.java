package system;

import java.util.ArrayList;
import java.util.List;

public class ParseHandler {

    // Helper to decode the "Length:Value" format
    public List<String> parsePacket(String packet) {
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

    public Integer tryParseInt(String value) {
        try {
            return Integer.valueOf(value);
        } catch (NumberFormatException e) {
            return null; // Instead of crashing, we return null to handle it gracefully
        }
    }
}
