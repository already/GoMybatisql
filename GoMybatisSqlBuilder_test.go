package GoMybatis

import (
	"fmt"
	"github.com/already/batisql/v7/engines"
	"github.com/already/batisql/v7/utils"
	"github.com/beevik/etree"
	"testing"
	"time"
)

//压力测试 sql构建情况
func Benchmark_SqlBuilder(b *testing.B) {
	b.StopTimer()
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <if test="name != nil">
            and name like concat('%',#{name},'%')
        </if>
        <if test="startTime != nil">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != nil">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>
    </select>
</mapper>`

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, false)

	var mapperTree = LoadMapperXml([]byte(mapper))
	var nodes = builder.nodeParser.Parser(mapperTree["selectByCondition"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	paramMap["name"] = "sssssssss"
	paramMap["startTime"] = time.Now()
	paramMap["endTime"] = time.Now().Add(time.Hour * 24)
	paramMap["page"] = 12
	paramMap["size"] = 2

	//paramMap["func_name != nil"] = func(arg map[string]interface{}) interface{} {
	//	return arg["name"] != nil
	//}
	//paramMap["func_startTime != nil"] = func(arg map[string]interface{}) interface{} {
	//	return arg["startTime"] != nil
	//}
	//paramMap["func_endTime != nil"] = func(arg map[string]interface{}) interface{} {
	//	return arg["endTime"] != nil
	//}
	paramMap["func_page >= 0 and size != 0"] = func(arg map[string]interface{}) interface{} {
		return arg["list"] != nil && arg["size"] != nil
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var array = []interface{}{}
		_, e := builder.BuildSql(paramMap, nodes, &array)
		if e != nil {
			b.Fatal(e)
		}
	}
}

//测试sql生成tps
func Test_SqlBuilder_Tps(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <if test="name != nil">
            and name like concat('%',#{name},'%')
        </if>
        <if test="startTime != nil">
            and create_time >= #{startTime}
        </if>
        <if test="endTime != nil">
            and create_time &lt;= #{endTime}
        </if>
        order by create_time desc
        <if test="page >= 0 and size != 0">limit #{page}, #{size}</if>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, false)
	var paramMap = make(map[string]interface{})
	paramMap["name"] = ""
	paramMap["startTime"] = ""
	paramMap["endTime"] = ""
	paramMap["page"] = 0
	paramMap["size"] = 0

	var nodes = builder.nodeParser.Parser(mapperTree["selectByCondition"].(*etree.Element).Child)

	var startTime = time.Now()
	for i := 0; i < 100000; i++ {
		//var sql, e =
		var array = []interface{}{}
		_, e := builder.BuildSql(paramMap, nodes, &array)
		if e != nil {
			t.Fatal(e)
		}
		//fmt.Println(sql, e)
	}
	utils.CountMethodTps(100000, startTime, "Test_SqlBuilder_Tps")
}

func TestGoMybatisSqlBuilder_BuildSqlForeach1(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
	<select id="selectPostIn" resultType="domain.blog.Post">
	select s_id,count(s_id) as mcount from drugproduct3
      <where>
       <if test="cityCode != nil">
          and city_code = #{cityCode}
       </if> 
       <!-- [ lon, lat] -->
       <if test="lon != nil and lat!=nil and distance!=nil">
          and ST_Distance(location,ST_WKTToSQL('POINT (#{lon} #{lat})')) > #{distance}
       </if>
       <if test="list != nil">
        and 
	   <foreach item="item" index="index" collection="list"
		open="(" separator="or" close=")">
		 (med_list_code=#{item.medListCode} and unit=#{item.unit} and stock> #{item.reqNum} )
		</foreach>
       </if>
       and sd_status =1 and dp_status =1 group by s_id having count(s_id) >= ${func_list_len}
	</select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, true)
	var nodes = builder.nodeParser.Parser(mapperTree["selectPostIn"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	paramMap["cityCode"] = "441800"
	paramMap["lon"] = -71.32
	paramMap["lat"] = 41.11
	paramMap["distance"] = 2100
	list := []map[string]interface{}{{"medListCode": "XA07FAK084N001010100041", "unit": "盒", "reqNum": 2},
		{"medListCode": "ZA04BAX0061020200423", "unit": "包", "reqNum": 1}}
	paramMap["list"] = list

	paramMap["func_list_len"] = func(arg map[string]interface{}) interface{} {
		fmt.Println("---------!!!!!--------->")
		return arg["list"]
	}

	var array = []interface{}{}

	var sql, err = builder.BuildSql(paramMap, nodes, &array)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql)
}

//
func TestGoMybatisSqlBuilder_BuildSqlForeach(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
	<select id="selectPostIn" resultType="domain.blog.Post">
	SELECT *
	FROM POST P
	WHERE ID in
		<foreach item="item" index="index" collection="list"
		open="(" separator="," close=")">
		#{item.id}
		</foreach>
	</select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, true)
	var nodes = builder.nodeParser.Parser(mapperTree["selectPostIn"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	list := []map[string]interface{}{{"id": "11333"}, {"id": "222"}}
	paramMap["list"] = list
	paramMap["startTime"] = time.Now()
	paramMap["endTime"] = time.Now().Add(time.Hour * 24)
	paramMap["page"] = 12
	paramMap["size"] = 2

	var array = []interface{}{}

	var sql, err = builder.BuildSql(paramMap, nodes, &array)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql)
}

func TestGoMybatisSqlBuilder_BuildSql(t *testing.T) {
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <resultMap id="BaseResultMap">
        <id column="id" property="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int"/>
    </resultMap>
    <select id="selectByCondition" resultMap="BaseResultMap">
        <bind name="pattern" value="'%' + name + '%'"/>
        select * from biz_activity
        <where>
            <if test="name != nil">
                and name like #{pattern}
            </if>
            <if test="startTime != nil">and create_time >= #{startTime}</if>
            <if test="endTime != nil">and create_time &lt;= #{endTime}</if>
        </where>
        order by 
        <trim prefix="" suffix="" suffixOverrides=",">
            <if test="name != nil">name,</if>
        </trim>
        desc
        <choose>
            <when test="page < 1">limit 3</when>
            <when test="page > 1">limit 2</when>
            <otherwise>limit 1</otherwise>
        </choose>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, true)
	var nodes = builder.nodeParser.Parser(mapperTree["selectByCondition"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	paramMap["name"] = "sssssssss"
	paramMap["startTime"] = time.Now()
	paramMap["endTime"] = time.Now().Add(time.Hour * 24)
	paramMap["page"] = 12
	paramMap["size"] = 2

	var array = []interface{}{}

	var sql, err = builder.BuildSql(paramMap, nodes, &array)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sql)
}

//压力测试 sql构建情况
func Benchmark_SqlBuilder_If_Element(b *testing.B) {
	b.StopTimer()
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
        <if test="name != nil">
        </if>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, false)
	var nodes = builder.nodeParser.Parser(mapperTree["selectByCondition"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	paramMap["name"] = ""
	paramMap["startTime"] = ""
	paramMap["endTime"] = ""
	paramMap["page"] = 0
	paramMap["size"] = 0

	//paramMap["type_name"] = StringType
	//paramMap["type_startTime"] = StringType
	//paramMap["type_endTime"] = StringType
	//paramMap["type_page"] = IntType
	//paramMap["type_size"] = IntType

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var array = []interface{}{}
		builder.BuildSql(paramMap, nodes, &array)
	}
}

//压力测试 element嵌套构建情况
func Benchmark_SqlBuilder_Nested(b *testing.B) {
	b.StopTimer()
	var mapper = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper>
    <!--List<Activity> selectByCondition(@Param("name") String name,@Param("startTime") Date startTime,@Param("endTime") Date endTime,@Param("index") Integer index,@Param("size") Integer size);-->
    <!-- 后台查询产品 -->
    <select id="selectByCondition">
        select * from biz_activity where delete_flag=1
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>
        <set>

        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
        </set>
    </select>
</mapper>`
	var mapperTree = LoadMapperXml([]byte(mapper))

	var builder = GoMybatisSqlBuilder{}.New(ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true), &LogStandard{}, false)
	var nodes = builder.nodeParser.Parser(mapperTree["selectByCondition"].(*etree.Element).Child)

	var paramMap = make(map[string]interface{})
	paramMap["name"] = ""
	paramMap["startTime"] = ""
	paramMap["endTime"] = ""
	paramMap["page"] = 0
	paramMap["size"] = 0

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var array = []interface{}{}
		_, e := builder.BuildSql(paramMap, nodes, &array)
		if e != nil {
			b.Fatal(e)
		}
	}
}
