package data

import (

"database/sql"

"geektime-go-orm/orm/db/register"

"geektime-go-orm/orm/predicate"

)




    const (
    
        BaseInfoDetail = "Detail"
    
        BaseInfoDescription = "Description"
    
    )

   
        
            func BaseInfoDetailLt(val string) predicate.Predicate {
                return predicate.C("Detail").Lt(predicate.Valuer{Value: val})
            }
        
            func BaseInfoDetailEq(val string) predicate.Predicate {
                return predicate.C("Detail").Eq(predicate.Valuer{Value: val})
            }
        
            func BaseInfoDetailGt(val string) predicate.Predicate {
                return predicate.C("Detail").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func BaseInfoDescriptionLt(val string) predicate.Predicate {
                return predicate.C("Description").Lt(predicate.Valuer{Value: val})
            }
        
            func BaseInfoDescriptionEq(val string) predicate.Predicate {
                return predicate.C("Description").Eq(predicate.Valuer{Value: val})
            }
        
            func BaseInfoDescriptionGt(val string) predicate.Predicate {
                return predicate.C("Description").Gt(predicate.Valuer{Value: val})
            }
        
    


    const (
    
        UserId = "Id"
    
        UserUsername = "Username"
    
        UserEmail = "Email"
    
        UserBirthdate = "Birthdate"
    
        UserIsActive = "IsActive"
    
        UserBaseInfo = "BaseInfo"
    
    )

   
        
            func UserIdLt(val int) predicate.Predicate {
                return predicate.C("Id").Lt(predicate.Valuer{Value: val})
            }
        
            func UserIdEq(val int) predicate.Predicate {
                return predicate.C("Id").Eq(predicate.Valuer{Value: val})
            }
        
            func UserIdGt(val int) predicate.Predicate {
                return predicate.C("Id").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func UserUsernameLt(val string) predicate.Predicate {
                return predicate.C("Username").Lt(predicate.Valuer{Value: val})
            }
        
            func UserUsernameEq(val string) predicate.Predicate {
                return predicate.C("Username").Eq(predicate.Valuer{Value: val})
            }
        
            func UserUsernameGt(val string) predicate.Predicate {
                return predicate.C("Username").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func UserEmailLt(val string) predicate.Predicate {
                return predicate.C("Email").Lt(predicate.Valuer{Value: val})
            }
        
            func UserEmailEq(val string) predicate.Predicate {
                return predicate.C("Email").Eq(predicate.Valuer{Value: val})
            }
        
            func UserEmailGt(val string) predicate.Predicate {
                return predicate.C("Email").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func UserBirthdateLt(val string) predicate.Predicate {
                return predicate.C("Birthdate").Lt(predicate.Valuer{Value: val})
            }
        
            func UserBirthdateEq(val string) predicate.Predicate {
                return predicate.C("Birthdate").Eq(predicate.Valuer{Value: val})
            }
        
            func UserBirthdateGt(val string) predicate.Predicate {
                return predicate.C("Birthdate").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func UserIsActiveLt(val bool) predicate.Predicate {
                return predicate.C("IsActive").Lt(predicate.Valuer{Value: val})
            }
        
            func UserIsActiveEq(val bool) predicate.Predicate {
                return predicate.C("IsActive").Eq(predicate.Valuer{Value: val})
            }
        
            func UserIsActiveGt(val bool) predicate.Predicate {
                return predicate.C("IsActive").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func UserBaseInfoLt(val JsonData[BaseInfo]) predicate.Predicate {
                return predicate.C("BaseInfo").Lt(predicate.Valuer{Value: val})
            }
        
            func UserBaseInfoEq(val JsonData[BaseInfo]) predicate.Predicate {
                return predicate.C("BaseInfo").Eq(predicate.Valuer{Value: val})
            }
        
            func UserBaseInfoGt(val JsonData[BaseInfo]) predicate.Predicate {
                return predicate.C("BaseInfo").Gt(predicate.Valuer{Value: val})
            }
        
    


    const (
    
        testa = "a"
    
        testb = "b"
    
        testc = "c"
    
        testd = "d"
    
    )

   
        
            func testaLt(val predicate.Column) predicate.Predicate {
                return predicate.C("a").Lt(predicate.Valuer{Value: val})
            }
        
            func testaEq(val predicate.Column) predicate.Predicate {
                return predicate.C("a").Eq(predicate.Valuer{Value: val})
            }
        
            func testaGt(val predicate.Column) predicate.Predicate {
                return predicate.C("a").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func testbLt(val register.TableName) predicate.Predicate {
                return predicate.C("b").Lt(predicate.Valuer{Value: val})
            }
        
            func testbEq(val register.TableName) predicate.Predicate {
                return predicate.C("b").Eq(predicate.Valuer{Value: val})
            }
        
            func testbGt(val register.TableName) predicate.Predicate {
                return predicate.C("b").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func testcLt(val sql.NullString) predicate.Predicate {
                return predicate.C("c").Lt(predicate.Valuer{Value: val})
            }
        
            func testcEq(val sql.NullString) predicate.Predicate {
                return predicate.C("c").Eq(predicate.Valuer{Value: val})
            }
        
            func testcGt(val sql.NullString) predicate.Predicate {
                return predicate.C("c").Gt(predicate.Valuer{Value: val})
            }
        
    
        
            func testdLt(val []byte) predicate.Predicate {
                return predicate.C("d").Lt(predicate.Valuer{Value: val})
            }
        
            func testdEq(val []byte) predicate.Predicate {
                return predicate.C("d").Eq(predicate.Valuer{Value: val})
            }
        
            func testdGt(val []byte) predicate.Predicate {
                return predicate.C("d").Gt(predicate.Valuer{Value: val})
            }
        
    



