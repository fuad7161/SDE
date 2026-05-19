public class Main {

    public static void main(String[] args) {
        Student s1 = new Student("John");
        Student s2 = new Student("John");

        System.out.println(s1 == s2);
        System.out.println(s1.equals(s2));
    }
}
