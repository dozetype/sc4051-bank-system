package ui;
import java.util.Scanner;

public class UI{
    private final Scanner scanner;

    public UI(){
        scanner = new Scanner(System.in);
    }

    public String inputString(){
        String input = scanner.nextLine();
        while(input.isEmpty()){ //Won't accept empty input
            input = scanner.nextLine();
        }
        return input;
    }

    public int inputInt(){
        System.out.print("Enter Number: ");
        int number;
        while (true) {
            String input = scanner.nextLine();
            try {
                number = Integer.parseInt(input);
                break; // valid input, break the loop
            } catch (NumberFormatException e) {
                System.out.println("That's not a valid integer. Try again.");
            }
        }

        return number;
    }
}