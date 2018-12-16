# Changelog

## 1.4: Interactive with dialog context - 15-12-2018

* When a question is about Lord Byron, and the database has two persons "Lord Byron" asks "which one"
* Dialog context to store the selected person by the user

## 1.3: Simple DB-pedia queries - 24-09-2017

* Domain - knowledge base mapping changes from 1:n to n:m
* Support for Sparql bases
* Intermediate results are logged
* Optimization phase using knowledge base statistics

## 1.2: Command-line app "nli" - 06-05-2017

* An executable application with "answer" and "suggest subcommands"
* Use an existing javascript autosuggest line editor (Tag-it!) and create an example web app
* Build an example application from a configuration file
* Rebuild of log as a proper dependency and with productions

## 1.1: Quantifier Scoping - 11-04-2017

* handle scoped questions
    * One sentence with ALL and 2 as quantifiers
    * One sentence where the right quantifer outscopes the left
* examples from relationships
* new: parse tree as new step

## 1: simple full-circle nli - 28-02-2017

* language: english
* question types: yes/no, who, which, how many
* second order predicates, aggregations
* proper nouns
* real database access (MySQL)
* a few simple questions
* simple natural language responses
* working example