#include <iostream>
#include <memory>
#include <string>

using namespace std;

// --------------------------------------------------
//  ACCOUNT BASE CLASS
// --------------------------------------------------
class Account {
public:
    virtual ~Account() = default;
    virtual string getAccountType() const = 0;
};

// --------------------------------------------------
//  PERSONAL ACCOUNT
// --------------------------------------------------
class PersonalAccount : public Account {
public:
    string fullName;
    int loyaltyPoints = 0;

    PersonalAccount(const string& name) : fullName(name) {}

    string getAccountType() const override {
        return "Personal";
    }
};

// --------------------------------------------------
//  COMPANY ACCOUNT
// --------------------------------------------------
class CompanyAccount : public Account {
public:
    string companyName;
    int employeeCount;
    double creditLimit = 0;

    CompanyAccount(string cname, int empCount)
        : companyName(cname), employeeCount(empCount) {}

    string getAccountType() const override {
        return "Company";
    }
};

// --------------------------------------------------
//  WORKFLOW BASE CLASS
// --------------------------------------------------
class Workflow {
public:
    virtual ~Workflow() = default;
    virtual void run(Account* acc) = 0;
};

// --------------------------------------------------
//  PERSONAL ACCOUNT WORKFLOW
// --------------------------------------------------
class PersonalWorkflow : public Workflow {
public:
    void run(Account* acc) override {
        auto* pAcc = dynamic_cast<PersonalAccount*>(acc);
        if (!pAcc) return;

        cout << "\n[Personal Workflow Started]\n";
        cout << "Full Name: " << pAcc->fullName << endl;

        pAcc->loyaltyPoints += 10;
        cout << "Loyalty points updated: " << pAcc->loyaltyPoints << endl;
    }
};

// --------------------------------------------------
//  COMPANY ACCOUNT WORKFLOW
// --------------------------------------------------
class CompanyWorkflow : public Workflow {
public:
    void run(Account* acc) override {
        auto* cAcc = dynamic_cast<CompanyAccount*>(acc);
        if (!cAcc) return;

        cout << "\n[Company Workflow Started]\n";
        cout << "Company Name: " << cAcc->companyName << endl;

        // Hesaplama örneği: çalışan sayısına göre kredi limiti
        cAcc->creditLimit = cAcc->employeeCount * 1500;
        cout << "Assigned credit limit: " << cAcc->creditLimit << endl;
    }
};

// --------------------------------------------------
//  WORKFLOW ENGINE
// --------------------------------------------------
class WorkflowEngine {
public:
    void runWorkflow(Account* acc) {
        if (acc->getAccountType() == "Personal") {
            PersonalWorkflow wf;
            wf.run(acc);
        } 
        else if (acc->getAccountType() == "Company") {
            CompanyWorkflow wf;
            wf.run(acc);
        }
    }
};

// --------------------------------------------------
//  MAIN
// --------------------------------------------------
int main() {
    WorkflowEngine engine;

    PersonalAccount pAcc("Ahmet Yılmaz");
    CompanyAccount cAcc("TechCorp", 120);

    engine.runWorkflow(&pAcc);
    engine.runWorkflow(&cAcc);

    return 0;
}

// --------------------------------------------------
//  SAMPLE CLASS
// --------------------------------------------------
class SampleClass {
public:
    void sampleMethod() {
        std::cout << "foo" << std::endl;
        std::cout << "bar" << std::endl;
        std::cout << "baz" << std::endl;
        std::cout << "qux" << std::endl;

        for (size_t i = 0; i < count; i++)
        {
            /**
             * Sample comment
             * Another line of comment
             * Yet another line of comment
             * @var i This is a variable in the loop
             * @return void
             * @throws None
             * @see SampleClass
             * @since 1.0
             */
        }
        
    }
};

// --------------------------------------------------
//  SAMPLE CLASS USAGE
// --------------------------------------------------
class SampleClassUsage {
public:
    void useSampleClass() {
        SampleClass sample;
        sample.sampleMethod();
    }
    private:
    const size_t count = 10;
    /**
     * SampleClassUsage comment
     * Another line of comment
     * Yet another line of comment
     * @param count This is a private member variable
     * @param sample An instance of SampleClass
     * @return void
     * @throws None
     */
};C++17